<template>
  <div>
    <form v-on:submit.prevent="createMeow">
      <div class="input-group">
        <input v-model.trim="meowBody" type="text" class="form-control" placeholder="What's happening?">
        <div class="input-group-append">
          <button class="btn btn-primary" type="submit">Meow</button>
        </div>
      </div>
    </form>

    <div class="mt-4">
      <Meow v-for="meow in meows" :key="meow.id" :meow="meow" />
    </div>
  </div>
</template>

<script>
import { mapState } from 'vuex';
import Meow from '@/components/Meow';

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
  components: {
    Meow,
  },
};
</script>
