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
        data:{},
        confs:[],
        result:{}
    },
    actions: {
        ShowLoading(context){
            this.state.loadingcount += 1;
            console.log(this.state.loadingcount)
            this.state.loading = true;
            console.log(this.state.loading)
        },
        HiddenLoading(context){
            this.state.loadingcount -= 1;
            console.log(this.state.loadingcount)
            if (this.state.loadingcount == 0) {
                this.state.loading = false;
            }
            console.log(this.state.loading)
        },
        EnableEdit(context){
            this.state.editdisable = false;
        },
        DisableEdit(context){
            this.state.editdisable = true;
        },
        GetServs(context,data) {
            this.state.menuvalue = data['env']
            data['env'] = this.state.menuvalue
            axios
                .get('/api/servs', {params:data})
                .then(response => {
                    this.state.data = response.data})
        },
        GetConfs(context) {
            axios
                .get('/api/confs')
                .then(response => {
                    this.state.confs = response.data.data})
        },
        Conf(context, data) {
            axios
                .post('/api/conf', data)
                .then(response => (this.state.data = response.data))
        },
        Pack(context, serv) {
            return axios.post('/release/pack', serv)
                .then(response => (this.state.result = response.data))
        },
        Trans(context, serv) {
            return axios.post('/release/trans', serv)
                .then(response => (this.state.result = response.data))
        },
        Post(context, serv) {
            return axios.post('/release/post', serv)
                .then(response => (this.state.result = response.data))
        }
    }
})
