import { get as getSetting, save as saveSetting, subscribe as subscribeSetting } from '@/repository/settings';

const useSettingsStore = defineStore(
  'settings',
  {
    state: () => ({
      inited: false,
      theme: { content: 'dark' },
      language: { content: 'zh' }
    }),
    actions: {
      switchToTheme(theme) {
        this.theme.content = theme
        saveSetting(this.theme).then(res => {
        })
      },
      async switchToLanguage(language) {
        if (language == this.language.content) {
          return
        } else {
          this.language.content = language
          await saveSetting(this.language);
        }
      },
      async getSettingInfo() {
        let theme = await getSetting('theme')
        let language = await getSetting('language')
        this.inited = true
        this.theme = theme
        this.language = language
      }
    }
  })

export default useSettingsStore
