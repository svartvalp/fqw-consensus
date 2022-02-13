<template>
  <div class="col">
  <div :class="['card', this.getBG, this.getText]">
    <div class="card-body">
      <h5 class="card-title">ID: {{this.state.local_id}}</h5>
      <div class="card-text" >Term: {{this.state.current_term}}</div>
      <div class="card-text">State: {{this.state.state}}</div>
      <div class="card-text">Host: {{this.state.host}}</div>
      <div class="card-text" v-if="this.state.voted_for">Voted for: {{this.state.voted_for}}</div>
      <div class="card-text" v-if="this.state.votes">Votes: {{this.state.votes}}</div>
      <div class="card-text" v-if="this.state.election_timeout && this.state.state !== 'Leader'">Election timeout: {{this.state.election_timeout}}</div>
    </div>
  </div>
  </div>

</template>
<script>

export default {
  name: 'State',
  components: {
  },
  setup() {
  },
  props: {
    state: Object
  },
  computed: {
    getBG() {
      if (this.state.state === 'Leader') {
        return 'bg-primary'
      }
      if (this.state.state === 'Dead') {
        return 'bg-warning'
      }

      if (this.state.state === 'Candidate') {
        return 'bg-info'
      }


      return "bg-secondary"
    },
    getText() {
      return "text-white"
    }
  }
}
</script>