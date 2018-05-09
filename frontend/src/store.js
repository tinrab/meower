import Vue from 'vue';
import Vuex from 'vuex';
import axios from 'axios';

const API_URL = 'http://localhost:8000';
const CREATE_MEOW_SUCCESS = 'CREATE_MEOW_SUCCESS';

Vue.use(Vuex);

export default new Vuex.Store({
  state: {
    meows: [],
  },
  mutations: {
    [CREATE_MEOW_SUCCESS](state, meow) {
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
      commit(CREATE_MEOW_SUCCESS, meow);
    },
  },
});
