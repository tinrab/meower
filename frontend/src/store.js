import Vue from 'vue';
import Vuex from 'vuex';
import axios from 'axios';
import VueNativeSock from 'vue-native-websocket';

const API_URL = 'http://localhost:8000';
const WS_URL = 'ws://localhost:8001/ws';
const CREATE_MEOW = 'CREATE_MEOW';

const MESSAGE_MEOW_CREATED = 1;

Vue.use(Vuex);

const store = new Vuex.Store({
  state: {
    meows: [],
  },
  mutations: {
    SOCKET_ONOPEN(state, event) {},
    SOCKET_ONCLOSE(state, event) {},
    SOCKET_ONERROR(state, event) {
      console.error(event);
    },
    SOCKET_ONMESSAGE(state, message) {
      switch (message.kind) {
        case MESSAGE_MEOW_CREATED:
          this.commit(CREATE_MEOW, { id: message.id, body: message.body });
      }
    },
    [CREATE_MEOW](state, meow) {
      state.meows = [meow, ...state.meows];
    },
  },
  actions: {
    async createMeow({ commit }, meow) {
      const res = await axios.post(`${API_URL}/meows`, null, {
        params: {
          body: meow.body,
        },
      });
      meow.id = res.data.id;
    },
  },
});

Vue.use(VueNativeSock, WS_URL, { store, format: 'json' });

export default store;
