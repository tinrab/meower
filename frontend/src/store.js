import Vue from 'vue';
import Vuex from 'vuex';
import axios from 'axios';
import VueNativeSock from 'vue-native-websocket';

const MEOW_URL = 'http://localhost:8080';
const QUERY_URL = 'http://localhost:8081';
const WS_URL = 'ws://localhost:8082/ws';

const SET_MEOWS = 'SET_MEOWS';
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
    [SET_MEOWS](state, meows) {
      state.meows = meows;
    },
    [CREATE_MEOW](state, meow) {
      state.meows = [meow, ...state.meows];
    },
  },
  actions: {
    getMeows({ commit }) {
      axios
        .get(`${QUERY_URL}/meows`)
        .then(({ data }) => {
          commit(SET_MEOWS, data);
        })
        .catch((err) => console.err(err));
    },
    async createMeow({ commit }, meow) {
      const { data } = await axios.post(`${MEOW_URL}/meows`, null, {
        params: {
          body: meow.body,
        },
      });
      meow.id = data.id;
    },
  },
});

Vue.use(VueNativeSock, WS_URL, { store, format: 'json' });

store.dispatch('getMeows');

export default store;
