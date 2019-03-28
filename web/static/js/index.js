import Vue from 'vue'
import {Button,Checkbox,CheckboxGroup,Col,Row,Menu,Submenu,MenuItem,Select,Option,Input,Loading,Tooltip} from 'element-ui';
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

axios.interceptors.request.use(function(config){  
    store.dispatch('ShowLoading')  
    return config  
},function(err){  
    return Promise.reject(err)  
});  
axios.interceptors.response.use(function(response){  
    store.dispatch('HiddenLoading')  
    return response  
},function(err){  
    return Promise.reject(err)  
}); 

new Vue({
  el: '#app',
  router,
  store,
  render: h => h(App)
})
