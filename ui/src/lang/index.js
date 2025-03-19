import { createI18n } from 'vue-i18n'
import en from './en'
import zh from './zh'

// 定义多语言配置
const messages = {
  en,  // 英文
  zh   // 中文
}

// 创建 i18n 实例
const i18n = createI18n({
  legacy: false, // 使用 Composition API，必须设置为 false
  locale: 'zh',  // 默认语言
  fallbackLocale: 'en', // 回退语言
  messages       // 语言包
})

export default i18n