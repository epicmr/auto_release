import Vue from 'vue'
import {Button,Checkbox,CheckboxGroup,Col,Row,Menu,Submenu,MenuItem,Select,Option,Input,Loading,Tooltip,Radio,RadioGroup,MessageBox,Message,Tag,Table,TableColumn,Form,FormItem,Pagination} from 'element-ui';
import 'element-ui/lib/theme-chalk/index.css';
import App from './App'
import router from './routes'
import store from './store'
import axios from 'axios'

Vue.config.debug = process.env.NODE_ENV !== 'production'
Vue.use(Button);
Vue.use(Checkbox);
Vue.use(CheckboxGroup);
Vue.use(Col);
Vue.use(Row);
Vue.use(Menu);
Vue.use(Submenu);
Vue.use(MenuItem);
Vue.use(Select);
Vue.use(Option);
Vue.use(Input);
Vue.use(Loading);
Vue.use(Tooltip);
Vue.use(Radio);
Vue.use(RadioGroup);
Vue.use(Tag);
Vue.use(Table);
Vue.use(TableColumn);
Vue.use(Form);
Vue.use(FormItem);
Vue.use(Pagination);
Vue.prototype.$msgbox = MessageBox;
Vue.prototype.$message = Message;

axios.interceptors.request.use(function(config){  
    store.commit('ShowLoading')  
    return config  
},function(err){  
    return promise.reject(err)  
});  
axios.interceptors.response.use((response) => {
    store.commit('HiddenLoading')  
        console.log("yyy")
    return response  
}, function (err) {
    if (302 === err.response.status) {
        console.log("xxx")
        window.location = '/login';
    } else {
        return promise.reject(err);
    }
});

new Vue({
  el: '#app',
  router,
  store,
  render: h => h(App)
})
