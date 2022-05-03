<template>
  <div class="col">
  <div :class="['card', this.getBG, this.getText, 'bg-opacity-50']">
    <div class="card-body">
      <h5 class="card-title">ID: {{this.state.local_id}} <router-link v-if="this.state.local_id" :to="'/servers/' + this.state.local_id">Перейти</router-link> </h5>
      <ul class="list-group list-group-flush">
        <li class="list-group-item" ><div class="fw-bold">Term:</div>  {{this.state.current_term}}</li>
        <li class="list-group-item"><div class="fw-bold">State:</div> <span :class="['badge',this.getBG]">{{this.state.state}}</span></li>
        <li class="list-group-item"><div class="fw-bold">Host:</div> {{this.state.host}}</li>
        <li class="list-group-item" v-if="this.state.voted_for"><div class="fw-bold">Voted for:</div> {{this.state.voted_for}}</li>
        <li class="list-group-item" v-if="this.state.votes"><div class="fw-bold">Votes:</div> {{this.state.votes}}</li>
        <li class="list-group-item" v-if="this.state.commit_index"><div class="fw-bold">Commit index:</div> {{this.state.commit_index}}</li>
        <li class="list-group-item" v-if="this.state.election_timeout && this.state.state !== 'Leader'"><div class="fw-bold">Election timeout:</div> {{this.getElectionTimeout}} ms</li>
        <li class="list-group-item" v-if="this.state.store"><div class="fw-bold">Store:</div> {{this.state.store}}</li>
      </ul>
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
    state: Object,
  },
  computed: {
    getBG() {
      if (this.state.state === 'Leader') {
        return 'bg-success'
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
    },
    getElectionTimeout() {
      console.log(typeof this.state.election_timeout)
      let electionTimeoutTime = Date.parse( this.state.election_timeout)
      return Math.round(Math.abs(Date.now() - electionTimeoutTime) )
    }
  }
}
</script>