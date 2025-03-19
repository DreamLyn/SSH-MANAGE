const useTagsViewStore = defineStore(
  'tags-view',
  {
    state: () => ({
      cachedViews: [],
      cachedViewMap: {},//记录Key和SSH的映射
    }),
    actions: {
      addCachedView(view) {
        this.cachedViews.push(view)
        this.cachedViewMap[view.key] = view.sshId;
      },
      delCachedView(view) {
        return new Promise(resolve => {
          const index = this.cachedViews.indexOf(view.name)
          index > -1 && this.cachedViews.splice(index, 1)
          resolve([...this.cachedViews])
        })
      },

      delAllCachedViews(view) {
        return new Promise(resolve => {
          this.cachedViews = []
          resolve([...this.cachedViews])
        })
      },
    }
  })

export default useTagsViewStore
