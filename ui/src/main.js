import { createApp } from 'vue'

import Cookies from 'js-cookie'

import ElementPlus from 'element-plus'
import 'element-plus/dist/index.css'
import locale from 'element-plus/es/locale/lang/zh-cn'
import 'element-plus/theme-chalk/dark/css-vars.css'

import '@/assets/styles/index.scss'; // global css

import App from './App'
import i18n from './lang' // 引入国际化配置
import directive from './directive'; // directive
import router from './router'
import store from './store'

// svg图标
import SvgIcon from '@/components/SvgIcon'
import elementIcons from '@/components/SvgIcon/svgicon'
import 'virtual:svg-icons-register'

import './permission'; // permission control

import { addDateRange, handleTree, parseTime, resetForm } from '@/utils/ssh-manage'

// 分页组件
import Pagination from '@/components/Pagination'
// 图片预览组件
import ImagePreview from "@/components/ImagePreview"

//自定义初始化
const app = createApp(App)

// 全局方法挂载
app.config.globalProperties.parseTime = parseTime
app.config.globalProperties.resetForm = resetForm
app.config.globalProperties.handleTree = handleTree
app.config.globalProperties.addDateRange = addDateRange

// 全局组件挂载
app.component('Pagination', Pagination)
app.component('ImagePreview', ImagePreview)

app.use(router)
app.use(store)
app.use(elementIcons)
app.component('svg-icon', SvgIcon)
app.use(i18n) // 注册 i18n

directive(app)

// 使用element-plus 并且设置全局的大小
app.use(ElementPlus, {
  locale: locale,
  // 支持 large、default、small
  size: Cookies.get('size') || 'default'
})

app.mount('#app')

