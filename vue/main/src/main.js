import Vue from 'vue'
import {Button,Checkbox,CheckboxGroup,Col,Row,Menu,Submenu,MenuItem,Select,Option,Input,Loading,Tooltip,Radio,RadioGroup,MessageBox,Message,Tag,Table,TableColumn,Form,FormItem,Pagination} from 'element-ui';
import 'element-ui/lib/theme-chalk/index.css';
import App from './App'
import router from './router'
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
Vue.prototype.$prompt = MessageBox.prompt;

//Vue.prototype.$loading = Loading.service;
//Vue.prototype.$msgbox = MessageBox;
//Vue.prototype.$alert = MessageBox.alert;
//Vue.prototype.$confirm = MessageBox.confirm;
//Vue.prototype.$prompt = MessageBox.prompt;
//Vue.prototype.$notify = Notification;
//Vue.prototype.$message = Message;


Date.prototype.format = function(fmt) { 
    var o = { 
        "M+" : this.getMonth()+1,                 //月份 
        "d+" : this.getDate(),                    //日 
        "h+" : this.getHours(),                   //小时 
        "m+" : this.getMinutes(),                 //分 
        "s+" : this.getSeconds(),                 //秒 
        "q+" : Math.floor((this.getMonth()+3)/3), //季度 
        "S"  : this.getMilliseconds()             //毫秒 
    }; 
    if(/(y+)/.test(fmt)) {
        fmt=fmt.replace(RegExp.$1, (this.getFullYear()+"").substr(4 - RegExp.$1.length)); 
    }
    for(var k in o) {
        if(new RegExp("("+ k +")").test(fmt)){
            fmt = fmt.replace(RegExp.$1, (RegExp.$1.length==1) ? (o[k]) : (("00"+ o[k]).substr((""+ o[k]).length)));
        }
    }
    return fmt; 
}  

process.on('unhandledRejection', (reason, p) => {
    console.log('Unhandled Rejection at:', p, 'reason:', reason);
});

axios.interceptors.request.use(function(config){  
    store.commit('ShowLoading')  
    return config  
},function(err){  
    console.log("request err")
    return promise.reject(err)  
});  
axios.interceptors.response.use((response) => {
    store.commit('HiddenLoading')  
    return response  
}, function (err) {
    if (302 === err.response.status) {
        window.location = '/login';
    } else {
        console.log("response err")
        return promise.reject(err);
    }
});

new Vue({
  el: '#app',
  router,
  store,
  render: h => h(App)
})
