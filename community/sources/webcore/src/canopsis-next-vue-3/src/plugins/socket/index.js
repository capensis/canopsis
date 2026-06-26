import Socket from './services/socket'

export default {
  install(app) {
    const socket = new Socket()

    app.config.globalProperties.$socket = socket
    app.provide('$socket', socket)
  },
}
