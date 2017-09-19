import Vue from 'vue'
import Vuex from 'vuex'
import VuePaginate from 'vue-paginate'
import App from './App.vue'

Vue.use(VuePaginate)

new Vue({
  el: '#app',
  render: h => h(App)
})
