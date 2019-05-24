import Vue from 'vue'
import VueRouter from 'vue-router'
import Home from './components/Home'
import Welcome from './components/Welcome'
import Release from './components/Release'
import Login from './components/Login'
import EnvConf from './components/conf/EnvConf'
import HostConf from './components/conf/HostConf'
import Conf from './components/conf/Conf'
import UserRepair from './components/repair/UserRepair'
Vue.use(VueRouter)

const routes = [
    {
        path: '/',
        component: Home,
        meta : {
            requireAuth: true, 
        },
        children:[
            {
                path: '',
                component: Welcome
            },
            {
                path: 'release-(local|test|stg|seta|setb|setc|sete)',
                component: Release
            },
            {
                path: 'conf-env',
                component: EnvConf
            },
            {
                path: 'conf-host',
                component: HostConf
            },
            {
                path: 'conf',
                component: Conf
            },
            {
                path: 'repair-user',
                component: UserRepair
            }
        ]
    },
    {
        path: '/login',
        component: Login
    }
]

const router = new VueRouter({
    routes,
    scrollBehavior (to, from, savedPosition) {
        if (savedPosition) {
            return savedPosition
        } else {
            return { x: 0, y: 0 }
        }
    }
    //beforeEach(to,from,next) {
    //    console.log('need login')
    //    if(to.meta.requireAuth){
    //        if (store.getters.isLogin){
    //            next();
    //            console.log('dont need login')
    //        }else {
    //            next({
    //                path : '/login',
    //                query : {redirect : to.fullPath}
    //            })
    //            console.log('need login')
    //        }
    //    }else {
    //        next()
    //    }
    //}
})

export default router
