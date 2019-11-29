import Vue from 'vue'
import Router from 'vue-router'
import Login from '../components/Login'
import Regist from '../components/Regist'
import Home from '../components/Home'
import Welcome from '../components/Welcome'
import Release from '../components/Release'
import EnvConf from '../components/conf/EnvConf'
import HostConf from '../components/conf/HostConf'
import RouteConf from '../components/conf/RouteConf'
import PrivConf from '../components/conf/PrivConf'
import GroupConf from '../components/conf/GroupConf'
import ServConf from '../components/conf/ServConf'
import UserRepair from '../components/repair/UserRepair'
import VipRepair from '../components/repair/VipRepair'
import CouponRepair from '../components/repair/CouponRepair'
import DealRepair from '../components/repair/DealRepair'

const originalPush = Router.prototype.push;
Router.prototype.push = function push(location) {
  return originalPush.call(this, location).catch(err => err)
};

Vue.use(Router)

const routes = [
    {
        path: '/',
        name: 'Welcome',
        component: Welcome,
        meta : {
            requireAuth: true, 
        },
        children:[
            {
                path: 'session',
                name: 'Session',
                component: Welcome,
                children:[
                    {
                        path: 'login',
                        name: 'Login',
                        component: Login
                    },
                    {
                        path: 'regist',
                        name: 'Regist',
                        component: Regist,
                    }
                ]
            },
            {
                path: '',
                name: 'Home',
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
                        path: 'conf-route',
                        component: RouteConf
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
                        path: 'conf-serv',
                        component: ServConf
                    },
                    {
                        path: 'conf-priv',
                        component: PrivConf
                    },
                    {
                        path: 'conf-group',
                        component: GroupConf
                    },
                    {
                        path: 'repair-user',
                        component: UserRepair
                    },
                    {
                        path: 'repair-vip',
                        component: VipRepair
                    },
                    {
                        path: 'repair-coupon',
                        component: CouponRepair
                    },
                    {
                        path: 'repair-deal',
                        component: DealRepair
                    }
                ]
            },
            { path: '*', redirect: { name: 'Home' }}
        ]
    }
]

export default new Router({
    mode: 'history',
    routes,
    scrollBehavior (to, from, savedPosition) {
        if (savedPosition) {
            return savedPosition
        } else {
            return { x: 0, y: 0 }
        }
    }
})

