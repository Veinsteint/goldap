import router from './router'
import store from './store'
import { Message } from 'element-ui'
import NProgress from 'nprogress' // progress bar
import 'nprogress/nprogress.css' // progress bar style
import { getToken } from '@/utils/auth' // get token from cookie
import getPageTitle from '@/utils/get-page-title'

NProgress.configure({ showSpinner: false }) // NProgress Configuration

const whiteList = ['/login', '/auth-redirect'] // no redirect whitelist 
router.beforeEach(async(to, from, next) => {
  // start progress bar
  NProgress.start()

  // set page title
  document.title = getPageTitle(to.meta.title)

  // determine whether the user has logged in 
  const hasToken = getToken()

  if (hasToken) {
    if (to.path === '/login') {
      // if is logged in, redirect to the home page 
      next({ path: '/' })
      NProgress.done() 
    } else {
      // determine whether the user has obtained his permission roles through getInfo
      const hasRoles = store.getters.roles && store.getters.roles.length > 0
      if (hasRoles) {
        next()
      } else {
        try {
          // get user info
          const userInfo = await store.dispatch('user/getInfo')
          const { ID, roles } = userInfo
          const userinfoForRoutes = { id: ID, roles: roles }

          // generate accessible routes map based on roles (MUST do this first!)
          const accessRoutes = await store.dispatch('permission/generateRoutes', userinfoForRoutes)
          accessRoutes.push({ path: '*', redirect: '/404', hidden: true })
          router.addRoutes(accessRoutes)

          // check if user is admin (role ID = 1)
          const isAdmin = userInfo.roles && userInfo.roles.some(role => role.ID === 1)

          // redirect non-admin users to profile page (except allowed routes)
          const allowedPaths = ['/profile/index', '/service/passwordService', '/help']
          if (!isAdmin && !allowedPaths.includes(to.path)) {
            next({ path: '/profile/index', replace: true })
            NProgress.done()
            return
          }

          // navigate to the target route
          next({ ...to, replace: true })
        } catch (error) {
          // remove token and go to login page to re-login

          await store.dispatch('user/resetToken')
          Message.error(error || 'Has Error')
          next(`/login?redirect=${to.path}`)
          NProgress.done()
        }
      }
    }
  } else {
    /* has no token*/

    if (whiteList.indexOf(to.path) !== -1) {
      next()
    } else if (to.path === '/changePassword' || to.path === '/password-service' || to.path === '/help') {
      next({ replace: true })
      // NProgress.done()
    } else {
      // other pages that do not have permission to access are redirected to the login page.
      next(`/login?redirect=${to.path}`)
      NProgress.done()
    }
  }
})

router.afterEach(() => {
  // finish progress bar
  NProgress.done()
})
