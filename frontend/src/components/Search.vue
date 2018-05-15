<template>
  <div>
    <input @keyup="searchMeows" v-model.trim="query" type="text" class="form-control" placeholder="Search...">
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
      query: '',
    };
  },
  computed: mapState({
    meows: (state) => state.searchResults,
  }),
  methods: {
    searchMeows() {
      if (this.query != this.lastQuery) {
        this.$store.dispatch('searchMeows', this.query);
        this.lastQuery = this.query;
      }
    },
  },
  components: {
    Meow,
  },
};
</script>
