<template>
  <div>
    <div v-if="this.getState">
      <div :class="['card',' bg-opacity-25','bg-primary']">
        <div class="card-body">
          <h5 class="card-title">ID: {{ this.getState.local_id }} </h5>
          <ul class="list-group list-group-flush">
            <li class="list-group-item">
              <div class="fw-bold">Term:</div>
              {{ this.getState.current_term }}
            </li>
            <li class="list-group-item"><div class="fw-bold">State:</div> <span :class="['badge',this.getBG]">{{this.getState.state}}</span></li>
            <li class="list-group-item">
              <div class="fw-bold">Host:</div>
              {{ this.getState.host }}
            </li>
            <li class="list-group-item" v-if="this.getState.store">
              <div class="fw-bold">Store:</div>
              {{ this.getState.store }}
            </li>
          </ul>
          <table class="table">
            <thead>
            <tr>
              <th scope="col">
                Index
              </th>
              <th scope="col">
                Term
              </th>
              <th scope="col">
                Key
              </th>
              <th scope="col">
                Value
              </th>
            </tr>
            </thead>
            <tbody>
            <tr v-for="log in this.getState.logs" :key="log.index">
              <td>
                {{log.index}}
              </td>
              <td>
                {{log.term}}
              </td>
              <td>
                {{log.data.key}}
              </td>
              <td>
                {{log.data.key}}
              </td>
            </tr>
            </tbody>
          </table>
        </div>
      </div>
    </div>
  </div>
</template>
<script>


export default {
  name: 'ServerInfo',
  components: {},
  props: {
    states: Array,
    id: String,
  },
  computed: {
    getState() {
      return this.states.filter(st => st.local_id == this.id).find(Boolean)
    },
    getBG() {
      let state = this.getState
      if (state.state === 'Leader') {
        return 'bg-success'
      }
      if (state.state === 'Dead') {
        return 'bg-warning'
      }

      if (state.state === 'Candidate') {
        return 'bg-info'
      }


      return "bg-secondary"
    },
  }
}
</script>