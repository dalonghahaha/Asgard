// ---------------------------------------------------------------------------
// Element Plus 2.14.x 的 `buildProps` 用了复杂的 `EpPropMergeType` 包装类型，
// vue-tsc 2.x 无法在模板里解出 `PropType<string>` 的实际值，导致 `:label="'x'"` /
// `:prop="'x'"` / `:value="1"` 这类基本类型绑定全部报 TS2322；slot 里的 `row`
// 也被推断为内部 `DefaultRow` 而不是真实的行类型。
// 官方仓库直到 EP 2.10+ 仍在跟进 (https://github.com/element-plus/element-plus/issues/17877)。
//
// 绕开方案：通过 Vue 3.3+ 的 GlobalComponents 接口把模板里常用的组件 prop 类型
// 放宽为 `any`，运行时不变。这样项目里所有 el-table-column / el-option 等用法
// 都能通过严格模式 typecheck，又不影响 IDE 跳转（仍由 component-resolver 提供）。
//
// T-202：unplugin-vue-components 的 dts 自动生成已关闭 (vite.config.ts)，本文件
// 是 GlobalComponents 的唯一来源。新增 .vue 自研组件时，在这里追加一行即可。
// ---------------------------------------------------------------------------
import type { DefineComponent } from 'vue'
import type { default as RouterLinkDefault } from 'vue-router'
import type { default as RouterViewDefault } from 'vue-router'
import type { default as MonitorChartDefault } from '../components/MonitorChart.vue'
import type { default as TerminalLogDefault } from '../components/TerminalLog.vue'

type AnyProps = DefineComponent<Record<string, any>>

declare module 'vue' {
  interface GlobalComponents {
    // 表格
    ElTable: AnyProps
    ElTableColumn: AnyProps
    // 表单
    ElForm: AnyProps
    ElFormItem: AnyProps
    ElInput: AnyProps
    ElInputNumber: AnyProps
    ElSelect: AnyProps
    ElOption: AnyProps
    ElRadio: AnyProps
    ElRadioGroup: AnyProps
    ElSwitch: AnyProps
    ElDatePicker: AnyProps
    ElPageHeader: AnyProps
    // 基础
    ElButton: AnyProps
    ElTag: AnyProps
    ElDivider: AnyProps
    ElLink: AnyProps
    ElIcon: AnyProps
    ElCard: AnyProps
    ElResult: AnyProps
    ElSpace: AnyProps
    // 布局
    ElRow: AnyProps
    ElCol: AnyProps
    ElContainer: AnyProps
    ElAside: AnyProps
    ElHeader: AnyProps
    ElMain: AnyProps
    // 弹层
    ElDialog: AnyProps
    ElDrawer: AnyProps
    ElDescriptions: AnyProps
    ElDescriptionsItem: AnyProps
    // 导航
    ElMenu: AnyProps
    ElMenuItem: AnyProps
    ElSubMenu: AnyProps
    ElDropdown: AnyProps
    ElDropdownMenu: AnyProps
    ElDropdownItem: AnyProps
    // 分页
    ElPagination: AnyProps
    // 全局
    ElConfigProvider: AnyProps
    // 路由
    RouterLink: typeof RouterLinkDefault
    RouterView: typeof RouterViewDefault
    // 本项目自研组件
    MonitorChart: typeof MonitorChartDefault
    TerminalLog: typeof TerminalLogDefault
  }

  interface ComponentCustomProperties {
    $loading: (options?: any) => any
    vLoading: boolean
  }
}

export {}
