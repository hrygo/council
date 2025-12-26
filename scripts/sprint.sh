#!/bin/bash

# ============================================================================
# Sprint 迭代执行脚本
# ============================================================================
# 用途: 自动化执行 Sprint 迭代相关的任务
# 使用: ./scripts/sprint.sh [command]
# ============================================================================

set -e

# 颜色定义
BOLD='\033[1m'
CYAN='\033[36m'
GREEN='\033[32m'
YELLOW='\033[33m'
RED='\033[31m'
RESET='\033[0m'

# 项目根目录
ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
cd "$ROOT_DIR"

# ============================================================================
# 辅助函数
# ============================================================================

print_header() {
    echo -e "${BOLD}${CYAN}$1${RESET}"
    echo "════════════════════════════════════════════════════════════════"
}

print_success() {
    echo -e "${GREEN}✅ $1${RESET}"
}

print_warning() {
    echo -e "${YELLOW}⚠️  $1${RESET}"
}

print_error() {
    echo -e "${RED}❌ $1${RESET}"
}

# ============================================================================
# 核心功能
# ============================================================================

# 显示 Sprint 状态
sprint_status() {
    print_header "📊 Sprint 进度报告"
    
    if [ -f "docs/reports/QUICK_REFERENCE.md" ]; then
        cat docs/reports/QUICK_REFERENCE.md | grep -A 10 "一分钟总览" || print_warning "无法读取进度概览"
    else
        print_error "进度报告文件不存在: docs/reports/QUICK_REFERENCE.md"
        exit 1
    fi
    
    echo ""
}

# 显示本周计划
sprint_plan() {
    print_header "📅 本周 Sprint 计划"
    
    if [ -f "docs/reports/QUICK_REFERENCE.md" ]; then
        echo -e "${BOLD}本周重点:${RESET}"
        cat docs/reports/QUICK_REFERENCE.md | grep -A 15 "本周重点" | head -20 || print_warning "无法读取本周计划"
        
        echo ""
        echo -e "${CYAN}详细计划: docs/plans/Sprint9_Execution_Plan.md${RESET}"
    else
        print_error "快速参考文件不存在"
        exit 1
    fi
    
    echo ""
}

# 显示下一步行动
sprint_next() {
    print_header "🚀 下一步行动"
    
    if [ -f "docs/reports/QUICK_REFERENCE.md" ]; then
        cat docs/reports/QUICK_REFERENCE.md | grep -A 20 "下一步行动" || print_warning "无法读取下一步行动"
    else
        print_error "快速参考文件不存在"
        exit 1
    fi
    
    echo ""
}

# 显示任务看板
sprint_board() {
    print_header "📋 任务看板"
    
    if [ -f "docs/reports/task_board.md" ]; then
        cat docs/reports/task_board.md | grep -A 40 "看板视图" | head -45 || print_warning "无法读取任务看板"
        
        echo ""
        echo -e "${CYAN}完整看板: docs/reports/task_board.md${RESET}"
    else
        print_error "任务看板文件不存在"
        exit 1
    fi
    
    echo ""
}

# 更新进度报告
sprint_update() {
    print_header "🔄 更新 Sprint 进度报告"
    
    print_warning "此功能需要手动更新文档或使用 AI 工具重新生成"
    echo "请运行: qoder 并要求更新进度报告"
    
    echo ""
}

# 验证开发计划
sprint_validate() {
    print_header "✅ 验证开发计划"
    
    if [ -f "scripts/validate_dev_plan.py" ]; then
        python3 scripts/validate_dev_plan.py docs/development_plan.md
        print_success "开发计划验证完成"
    else
        print_warning "验证脚本不存在: scripts/validate_dev_plan.py"
    fi
    
    echo ""
}

# 完整的 Sprint 概览
sprint_overview() {
    print_header "🎯 Sprint 完整概览"
    echo ""
    
    sprint_status
    sprint_plan
    sprint_next
    
    echo -e "${BOLD}${GREEN}════════════════════════════════════════════════════════════════${RESET}"
    echo -e "${BOLD}${GREEN}📊 Sprint 概览完成${RESET}"
    echo -e "${BOLD}${GREEN}════════════════════════════════════════════════════════════════${RESET}"
    echo ""
}

# 检查 Sprint 健康度
sprint_health() {
    print_header "🏥 Sprint 健康度检查"
    echo ""
    
    local issues=0
    
    # 检查必需文件
    echo -e "${BOLD}检查必需文件...${RESET}"
    
    files=(
        "docs/development_plan.md"
        "docs/reports/QUICK_REFERENCE.md"
        "docs/reports/task_board.md"
        "docs/reports/task_progress_report.md"
        "docs/plans/Sprint9_Execution_Plan.md"
    )
    
    for file in "${files[@]}"; do
        if [ -f "$file" ]; then
            echo -e "  ${GREEN}✓${RESET} $file"
        else
            echo -e "  ${RED}✗${RESET} $file (缺失)"
            ((issues++))
        fi
    done
    
    echo ""
    
    # 检查文档时效性 (最后修改时间)
    echo -e "${BOLD}检查文档时效性...${RESET}"
    
    if [ -f "docs/reports/QUICK_REFERENCE.md" ]; then
        last_update=$(stat -f "%Sm" -t "%Y-%m-%d" docs/reports/QUICK_REFERENCE.md 2>/dev/null || stat -c "%y" docs/reports/QUICK_REFERENCE.md 2>/dev/null | cut -d' ' -f1)
        today=$(date +%Y-%m-%d)
        
        if [ "$last_update" = "$today" ]; then
            echo -e "  ${GREEN}✓${RESET} 快速参考今日已更新"
        else
            echo -e "  ${YELLOW}⚠${RESET} 快速参考上次更新: $last_update (建议每日更新)"
            ((issues++))
        fi
    fi
    
    echo ""
    
    # 总结
    if [ $issues -eq 0 ]; then
        print_success "Sprint 健康度良好！"
    else
        print_warning "发现 $issues 个问题，建议修复"
    fi
    
    echo ""
}

# ============================================================================
# 命令行参数处理
# ============================================================================

show_help() {
    echo ""
    echo -e "${BOLD}${CYAN}Sprint 迭代执行脚本${RESET}"
    echo "════════════════════════════════════════════════════════════════"
    echo ""
    echo "用途: 管理和追踪 Sprint 迭代进度"
    echo ""
    echo -e "${BOLD}使用方法:${RESET}"
    echo "  ./scripts/sprint.sh [command]"
    echo ""
    echo -e "${BOLD}可用命令:${RESET}"
    echo "  status      - 显示当前 Sprint 状态"
    echo "  plan        - 显示本周 Sprint 计划"
    echo "  next        - 显示下一步行动"
    echo "  board       - 显示任务看板"
    echo "  overview    - 显示完整 Sprint 概览"
    echo "  health      - 检查 Sprint 健康度"
    echo "  update      - 更新进度报告 (需手动或 AI)"
    echo "  validate    - 验证开发计划"
    echo "  help        - 显示此帮助信息"
    echo ""
    echo -e "${BOLD}示例:${RESET}"
    echo "  ./scripts/sprint.sh status"
    echo "  ./scripts/sprint.sh overview"
    echo "  ./scripts/sprint.sh health"
    echo ""
}

# 主程序
main() {
    case "${1:-help}" in
        status)
            sprint_status
            ;;
        plan)
            sprint_plan
            ;;
        next)
            sprint_next
            ;;
        board)
            sprint_board
            ;;
        overview)
            sprint_overview
            ;;
        update)
            sprint_update
            ;;
        validate)
            sprint_validate
            ;;
        health)
            sprint_health
            ;;
        help|--help|-h)
            show_help
            ;;
        *)
            print_error "未知命令: $1"
            show_help
            exit 1
            ;;
    esac
}

main "$@"
