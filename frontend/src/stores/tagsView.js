import { defineStore } from 'pinia'

export const useTagsViewStore = defineStore('tagsView', {
    state: () => ({
        visitedViews: [],
        cachedViews: []
    }),
    actions: {
        addView(view) {
            this.addVisitedView(view)
            this.addCachedView(view)
        },
        addVisitedView(view) {
            if (this.visitedViews.some(v => v.path === view.path)) return
            this.visitedViews.push(Object.assign({}, view, {
                title: view.meta.title || 'no-name'
            }))
        },
        addCachedView(view) {
            if (this.cachedViews.includes(view.name)) return
            if (!view.meta.noCache) {
                this.cachedViews.push(view.name)
            }
        },
        delView(view) {
            return new Promise(resolve => {
                this.delVisitedView(view)
                this.delCachedView(view)
                resolve({
                    visitedViews: [...this.visitedViews],
                    cachedViews: [...this.cachedViews]
                })
            })
        },
        delVisitedView(view) {
            for (const [i, v] of this.visitedViews.entries()) {
                if (v.path === view.path) {
                    this.visitedViews.splice(i, 1)
                    break
                }
            }
        },
        delCachedView(view) {
            const index = this.cachedViews.indexOf(view.name)
            index > -1 && this.cachedViews.splice(index, 1)
        },
        delOthersViews(view) {
            return new Promise(resolve => {
                this.delOthersVisitedViews(view)
                this.delOthersCachedViews(view)
                resolve({
                    visitedViews: [...this.visitedViews],
                    cachedViews: [...this.cachedViews]
                })
            })
        },
        delOthersVisitedViews(view) {
            this.visitedViews = this.visitedViews.filter(v => {
                return v.meta.affix || v.path === view.path
            })
        },
        delOthersCachedViews(view) {
            const index = this.cachedViews.indexOf(view.name)
            if (index > -1) {
                this.cachedViews = this.cachedViews.slice(index, index + 1)
            } else {
                this.cachedViews = []
            }
        },
        delAllViews() {
            return new Promise(resolve => {
                this.delAllVisitedViews()
                this.delAllCachedViews()
                resolve({
                    visitedViews: [...this.visitedViews],
                    cachedViews: [...this.cachedViews]
                })
            })
        },
        delAllVisitedViews() {
            // keep affix tags
            const affixTags = this.visitedViews.filter(tag => tag.meta.affix)
            this.visitedViews = affixTags
        },
        delAllCachedViews() {
            this.cachedViews = []
        },
        updateVisitedView(view) {
            for (let v of this.visitedViews) {
                if (v.path === view.path) {
                    v = Object.assign(v, view)
                    break
                }
            }
        }
    }
})
