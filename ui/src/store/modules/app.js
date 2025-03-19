import Cookies from 'js-cookie'

const useAppStore = defineStore(
  'app',
  {
    state: () => ({   
      device: 'desktop',
      size: Cookies.get('size') || 'default'
    }),
    actions: {
      toggleDevice(device) {
        this.device = device
      },
      setSize(size) {
        this.size = size;
        Cookies.set('size', size)
      }
    }
  })

export default useAppStore
