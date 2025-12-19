# SPEC-501: i18n 国际化完整实现

> **优先级**: P0  
> **类型**: Feature  
> **预估工时**: 8h

## 1. 概述

实现完整的多语言国际化支持，允许用户在 EN/ZH/JP 之间切换界面语言。

## 2. 目标

- 支持 2 种语言: English (en), 简体中文 (zh-CN)
- 语言切换即时生效，无需刷新页面
- 语言偏好持久化到 localStorage
- 自动检测浏览器语言作为默认值

## 3. 技术方案

### 3.1 依赖库

```json
{
  "i18next": "^25.x",
  "react-i18next": "^16.x",
  "i18next-browser-languagedetector": "^8.x"
}
```

### 3.2 目录结构

```
frontend/src/
├── i18n/
│   ├── index.ts              # i18n 初始化配置
│   ├── locales/
│   │   ├── en/
│   │   │   ├── common.json   # 通用文案
│   │   │   ├── workflow.json # 工作流相关
│   │   │   └── chat.json     # 聊天相关
│   │   ├── zh-CN/
│   │   │   ├── common.json
│   │   │   ├── workflow.json
│   │   │   └── chat.json
```

### 3.3 组件改造

**Before:**
```tsx
<button>Save Workflow</button>
```

**After:**
```tsx
import { useTranslation } from 'react-i18next';

function SaveButton() {
  const { t } = useTranslation('workflow');
  return <button>{t('actions.save')}</button>;
}
```

### 3.4 语言切换器组件

```tsx
// components/LanguageSwitcher.tsx
export function LanguageSwitcher() {
  const { i18n } = useTranslation();
  
  const languages = [
    { code: 'en', label: 'EN' },
    { code: 'zh-CN', label: '中文' },
  ];

  return (
    <select 
      value={i18n.language} 
      onChange={(e) => i18n.changeLanguage(e.target.value)}
    >
      {languages.map(({ code, label }) => (
        <option key={code} value={code}>{label}</option>
      ))}
    </select>
  );
}
```

## 4. 待翻译内容清单

### 4.1 通用 (common.json)

| Key              | EN         | ZH-CN     |
| :--------------- | :--------- | :-------- |
| `actions.save`   | Save       | 保存      |
| `actions.cancel` | Cancel     | 取消      |
| `actions.delete` | Delete     | 删除      |
| `actions.edit`   | Edit       | 编辑      |
| `actions.create` | Create     | 创建      |
| `status.loading` | Loading... | 加载中... |
| `status.error`   | Error      | 错误      |
| `status.success` | Success    | 成功      |

### 4.2 工作流 (workflow.json)

| Key                 | EN               | ZH-CN        |
| :------------------ | :--------------- | :----------- |
| `nodes.agent`       | Agent            | 智能体       |
| `nodes.vote`        | Vote             | 投票         |
| `nodes.loop`        | Loop             | 循环         |
| `nodes.factCheck`   | Fact Check       | 事实核查     |
| `nodes.humanReview` | Human Review     | 人工审核     |
| `builder.title`     | Workflow Builder | 工作流编辑器 |
| `estimate.title`    | Cost Estimation  | 成本预估     |

### 4.3 聊天 (chat.json)

| Key                 | EN                | ZH-CN       |
| :------------------ | :---------------- | :---------- |
| `input.placeholder` | Type a message... | 输入消息... |
| `panel.title`       | Chat              | 对话        |

## 5. 验收标准

- [ ] 语言切换器显示在 Header 右侧
- [ ] 切换语言后，所有 UI 文案即时更新
- [ ] 刷新页面后，语言偏好保持
- [ ] 无硬编码字符串残留 (lint 规则检查)
- [ ] 所有两种语言翻译完整率 100%

## 6. 测试用例

```typescript
describe('i18n', () => {
  it('should detect browser language', () => {
    // Mock navigator.language = 'zh-CN'
    expect(i18n.language).toBe('zh-CN');
  });

  it('should switch language', async () => {
    await i18n.changeLanguage('ja');
    expect(screen.getByText('保存')).not.toBeInTheDocument();
    expect(screen.getByText('保存')).toBeInTheDocument();
  });

  it('should persist language preference', () => {
    localStorage.setItem('i18nextLng', 'en');
    // Reload and check
    expect(i18n.language).toBe('en');
  });
});
```

## 7. 风险与缓解

| 风险              | 缓解措施                             |
| :---------------- | :----------------------------------- |
| 翻译质量参差不齐  | 使用 AI 初译 + 人工审核              |
| 遗漏硬编码字符串  | 添加 ESLint 规则 `no-literal-string` |
| 日期/数字格式差异 | 使用 `Intl` API 进行格式化           |
