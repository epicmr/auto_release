import Vue from 'vue'
import Vuex from 'vuex'
import axios from 'axios'

Vue.prototype.$http = axios
Vue.use(Vuex)

const debug = process.env.NODE_ENV !== 'production'

export default new Vuex.Store({
    state: {
        isLogin: false,
        loadingcount: 0,
        loading: false,
        editdisable: true,
        menuvalue: "",
        envs:[],
        hosts:[],
        confs:[],
        data:{}
    },
    mutations: {
        ShowLoading(state){
            state.loadingcount += 1;
            console.log(state.loadingcount)
            state.loading = true;
            console.log(state.loading)
        },
        HiddenLoading(state){
            state.loadingcount -= 1;
            console.log(state.loadingcount)
            if (state.loadingcount == 0) {
                state.loading = false;
            }
            console.log(state.loading)
        },
        EnableEdit(state){
            state.editdisable = false;
        },
        DisableEdit(state){
            state.editdisable = true;
        },
        GetServs(state,data) {
            state.menuvalue = data['env']
            data['env'] = state.menuvalue
            axios
                .get('/api/servs', {params:data})
                .then(response => {
                    state.data = response.data.data})
        },
        GetEnvs(state,data) {
            axios
                .get('/api/envs', {params:data})
                .then(response => {
                    state.envs = response.data.data})
        },
        Env(state, data) {
            data.hosts = null
            axios
                .post('/api/env', data)
                .then(response => (state.data = response.data))
        },
        GetHosts(state) {
            axios
                .get('/api/hosts')
                .then(response => {
                    state.hosts = response.data.data})
        },
        Host(state, data) {
            console.log(data)
            axios
                .post('/api/host', data)
                .then(response => (state.data = response.data))
        },
        GetConfs(state) {
            axios
                .get('/api/confs')
                .then(response => {
                    state.confs = response.data.data})
        },
        Conf(state, data) {
            axios
                .post('/api/conf', data)
                .then(response => (state.data = response.data))
        },
        Pack(state, serv) {
            return axios.post('/release/pack', serv)
                .then(response => (state.data = response.data))
        },
        Trans(state, serv) {
            return axios.post('/release/trans', serv)
                .then(response => (state.data = response.data))
        },
        Post(state, serv) {
            return axios.post('/release/post', serv)
                .then(response => (state.data = response.data))
        }
    }
})
