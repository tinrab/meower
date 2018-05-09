<template>
  <div>
    <h1>Meower</h1>
    <div>
      <input v-model.trim="meowBody">
      <button @click="createMeow">Meow</button>
    </div>
    <div>
      <ul>
        <li v-for="meow in meows" :key="meow.id">{{meow.body}}</li>
      </ul>
    </div>
  </div>
</template>

<script>
import { mapState } from 'vuex';

export default {
  data() {
    return {
      meowBody: '',
    };
  },
  computed: mapState({
    meows: (state) => state.meows,
  }),
  methods: {
    createMeow() {
      if (this.meowBody.length != 0) {
        this.$store.dispatch('createMeow', { body: this.meowBody });
        this.meowBody = '';
      }
    },
  },
};
</script>
