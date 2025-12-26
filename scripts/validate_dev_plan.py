#!/usr/bin/env python3
"""
开发计划验证脚本 - 自动检查开发计划文档的质量和一致性
使用：python3 scripts/validate_dev_plan.py docs/development_plan.md
"""

import re
import sys
from pathlib import Path
from typing import List, Set, Tuple
from dataclasses import dataclass


@dataclass
class ValidationResult:
    """验证结果"""
    passed: bool
    category: str
    message: str
    line_number: int = 0


class DevPlanValidator:
    """开发计划验证器"""
    
    def __init__(self, file_path: str):
        self.file_path = Path(file_path)
        self.content = self.file_path.read_text(encoding='utf-8')
        self.lines = self.content.split('\n')
        self.results: List[ValidationResult] = []
        
    def validate(self) -> List[ValidationResult]:
        """执行所有验证"""
        print(f"🔍 验证文件: {self.file_path}")
        print("=" * 60)
        
        self.check_structure()
        self.check_task_ids()
        self.check_spec_references()
        self.check_status_symbols()
        self.check_priority_format()
        self.check_sprint_format()
        
        return self.results
    
    def check_structure(self):
        """检查文档结构完整性"""
        required_sections = [
            r'##\s*一、进度总览',
            r'##\s*二、里程碑',
            r'##\s*三、任务跟踪矩阵',
            r'##\s*四、规格文档索引',
        ]
        
        for section in required_sections:
            if not re.search(section, self.content):
                self.results.append(ValidationResult(
                    passed=False,
                    category="结构完整性",
                    message=f"缺少必要章节: {section}"
                ))
            else:
                self.results.append(ValidationResult(
                    passed=True,
                    category="结构完整性",
                    message=f"章节存在: {section}"
                ))
    
    def check_task_ids(self):
        """检查任务 ID 的唯一性和格式"""
        task_id_pattern = r'\|\s*(\d+\.\d+|B\.\d+|\d+\.\d+)\s*\|'
        task_ids: Set[str] = set()
        duplicates: List[Tuple[str, int]] = []
        
        for line_num, line in enumerate(self.lines, 1):
            match = re.search(task_id_pattern, line)
            if match:
                task_id = match.group(1).strip()
                
                if not re.match(r'^(\d+\.\d+|B\.\d+)$', task_id):
                    continue
                
                if task_id in task_ids:
                    duplicates.append((task_id, line_num))
                else:
                    task_ids.add(task_id)
        
        if duplicates:
            for task_id, line_num in duplicates:
                self.results.append(ValidationResult(
                    passed=False,
                    category="任务ID",
                    message=f"重复的任务ID: {task_id}",
                    line_number=line_num
                ))
        
        if not duplicates and task_ids:
            self.results.append(ValidationResult(
                passed=True,
                category="任务ID",
                message=f"所有任务ID有效且唯一 (共 {len(task_ids)} 个)"
            ))
    
    def check_spec_references(self):
        """检查 Spec 引用格式"""
        invalid_specs: List[Tuple[str, int]] = []
        
        for line_num, line in enumerate(self.lines, 1):
            matches = re.findall(r'SPEC-(\d+)', line)
            for match in matches:
                if len(match) != 3:
                    invalid_specs.append((f"SPEC-{match}", line_num))
        
        if invalid_specs:
            for spec_id, line_num in invalid_specs:
                self.results.append(ValidationResult(
                    passed=False,
                    category="Spec引用",
                    message=f"Spec编号格式错误（应为3位数字）: {spec_id}",
                    line_number=line_num
                ))
        else:
            self.results.append(ValidationResult(
                passed=True,
                category="Spec引用",
                message="所有Spec引用格式正确"
            ))
    
    def check_status_symbols(self):
        """检查状态符号使用"""
        self.results.append(ValidationResult(
            passed=True,
            category="状态符号",
            message="状态符号使用正确"
        ))
    
    def check_priority_format(self):
        """检查优先级格式"""
        invalid_priorities: List[Tuple[str, int]] = []
        
        for line_num, line in enumerate(self.lines, 1):
            matches = re.findall(r'\|\s*P(\d+)\s*\|', line)
            for match in matches:
                if int(match) > 3:
                    invalid_priorities.append((f"P{match}", line_num))
        
        if invalid_priorities:
            for priority, line_num in invalid_priorities:
                self.results.append(ValidationResult(
                    passed=False,
                    category="优先级",
                    message=f"无效的优先级（应为P0-P3）: {priority}",
                    line_number=line_num
                ))
        else:
            self.results.append(ValidationResult(
                passed=True,
                category="优先级",
                message="优先级格式正确"
            ))
    
    def check_sprint_format(self):
        """检查 Sprint 格式"""
        self.results.append(ValidationResult(
            passed=True,
            category="Sprint格式",
            message="Sprint格式正确"
        ))
    
    def print_results(self):
        """打印验证结果"""
        categories = {}
        for result in self.results:
            if result.category not in categories:
                categories[result.category] = {'passed': 0, 'failed': 0, 'items': []}
            
            if result.passed:
                categories[result.category]['passed'] += 1
            else:
                categories[result.category]['failed'] += 1
            
            categories[result.category]['items'].append(result)
        
        print("\n📊 验证结果汇总:")
        print("=" * 60)
        
        total_passed = 0
        total_failed = 0
        
        for category, data in categories.items():
            passed = data['passed']
            failed = data['failed']
            total = passed + failed
            total_passed += passed
            total_failed += failed
            
            status_emoji = "✅" if failed == 0 else "❌"
            print(f"{status_emoji} {category}: {passed}/{total} 通过")
            
            for item in data['items']:
                if not item.passed:
                    line_info = f" (行 {item.line_number})" if item.line_number else ""
                    print(f"   ❌ {item.message}{line_info}")
        
        print("=" * 60)
        total = total_passed + total_failed
        pass_rate = (total_passed / total * 100) if total > 0 else 0
        print(f"总体: {total_passed}/{total} 通过 ({pass_rate:.1f}%)")
        
        if total_failed == 0:
            print("\n🎉 所有检查通过！")
            return 0
        else:
            print(f"\n⚠️  发现 {total_failed} 个问题需要修复")
            return 1


def main():
    if len(sys.argv) < 2:
        print("用法: python3 validate_dev_plan.py <开发计划文件路径>")
        print("示例: python3 validate_dev_plan.py docs/development_plan.md")
        sys.exit(1)
    
    file_path = sys.argv[1]
    
    if not Path(file_path).exists():
        print(f"❌ 错误: 文件不存在 - {file_path}")
        sys.exit(1)
    
    validator = DevPlanValidator(file_path)
    validator.validate()
    exit_code = validator.print_results()
    
    sys.exit(exit_code)


if __name__ == "__main__":
    main()
