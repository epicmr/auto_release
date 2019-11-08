import Vue from 'vue'
import Vuex from 'vuex'
import axios from 'axios'

var instance = axios
if (debug) {
    instance = axios.create({
        baseURL: 'http://beta-cgi.gstyun.cn:5480/',
        timeout: 3000,
        headers: {'content-type': 'application/json'}
    });
}

Vue.prototype.$http = instance
Vue.use(Vuex)

const debug = process.env.NODE_ENV !== 'production'

function handleErr(err) {
    console.log(err)
    if (debug) {
        console.info("%c [axiso catch error]", "color:orange", err)
    }
}

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
        servs:[],
        items:[],
        itemstree:[],
        data:{}
    },
    mutations: {
        ShowLoading(state){
            state.loadingcount += 1;
            state.loading = true;
        },
        HiddenLoading(state){
            state.loadingcount -= 1;
            if (state.loadingcount == 0) {
                state.loading = false;
            }
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
            instance
                .get('/api/servs', {params:data})
                .then(response => {state.servs = response.data.data})
                .catch(handleErr)
        },
        GetAllItems(state,data) {
            instance
                .get('/api/items', {params:data})
                .then(response => {state.items = response.data.data})
                .catch(handleErr)
        },
        GetItemsTree(state,data) {
            instance
                .get('/api/itemstree', {params:data})
                .then(response => {state.itemstree = response.data.data})
                .catch(handleErr)
        },
        Item(state, data) {
            instance
                .post('/api/item', data)
                .then(response => (state.data = response.data))
                .catch(handleErr)
        },
        GetEnvs(state,data) {
            instance
                .get('/api/envs', {params:data})
                .then(response => {state.envs = response.data.data})
                .catch(handleErr)
        },
        Env(state, data) {
            data.hosts = null
            instance
                .post('/api/env', data)
                .then(response => (state.data = response.data))
                .catch(handleErr)
        },
        UserGroup(state, data) {
            instance
                .post('/api/usergroup', data)
                .then(response => (state.data = response.data))
                .catch(handleErr)
        },
        GetHosts(state) {
            instance
                .get('/api/hosts')
                .then(response => {state.hosts = response.data.data})
                .catch(handleErr)
        },
        Host(state, data) {
            console.log(data)
            instance
                .post('/api/host', data)
                .then(response => (state.data = response.data))
                .catch(handleErr)
        },
        GetConfs(state) {
            instance
                .get('/api/confs')
                .then(response => {state.confs = response.data.data})
                .catch(handleErr)
        },
        Conf(state, data) {
            instance
                .post('/api/conf', data)
                .then(response => (state.data = response.data))
                .catch(handleErr)
        },
        Pack(state, serv) {
            return instance.post('/release/pack', serv)
                .then(response => {state.data = response.data})
                .catch(handleErr)
        },
        Trans(state, serv) {
            return instance.post('/release/trans', serv)
                .then(response => (state.data = response.data))
                .catch(handleErr)
        },
        Post(state, serv) {
            return instance.post('/release/post', serv)
                .then(response => (state.data = response.data))
                .catch(handleErr)
        },
        DeleteUser(state, data) {
            return instance.get('/third/deleteUser', {params:data})
                .then(response => (state.data = response.data))
                .catch(handleErr)
        }
    }
})
