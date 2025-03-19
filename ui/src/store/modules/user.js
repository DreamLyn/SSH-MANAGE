
import { authWithPassword, getAuthStore } from '@/repository/admin'

const useUserStore = defineStore(
  'user',
  {
    state: () => ({
      email: '',
    }),
    actions: {
      // 登录
      login(userInfo) {
        const username = userInfo.username.trim()
        const password = userInfo.password
        return new Promise((resolve, reject) => {
          (async () => {  // 立即执行的异步函数
            try {
              await authWithPassword(username.trim(), password);
              resolve();
            } catch (err) {
              reject(err);  // 使用 catch 中捕获的 err
            }
          })();
        })
      },
      // 退出系统
      logOut() {
        return new Promise((resolve, reject) => {
          const auth = getAuthStore()
          auth.clear()
          resolve()
        })
      }
    }
  })

export default useUserStore
