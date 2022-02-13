<template>
  <div class="container" style="height: 100%">
    <div class="row align-items-center h-100">
        <State v-for="state in states" :state="state" :key="state"></State>
    </div>
  </div>
</template>

<script>
import State from "./components/State";
import axios from 'axios';
import {reactive,ref} from "vue";
export default {
  name: 'App',
  components: {
    State,
  },
  setup() {
    let hosts = [
      "localhost:8081",
      "localhost:8082",
      "localhost:8083"
    ]
    const states = reactive(ref([]))

    const set = (st) => {
      states.value = st
    }

    const loadData = async (hosts)  => {
      setInterval(async function (){
        let newStates = []
        for (let i = 0; i < hosts.length; i++) {
          try {
            let res = await axios.get("http://" + hosts[i] + "/state")
            let state = res.data
            state.host = hosts[i]
            newStates.push(state)
          } catch (e)  {
            newStates.push({
              state: 'Dead',
              host: hosts[i]
            })
            // console.log(e)
          }
        }
        set(newStates)
      }, 50)
    }

    loadData(hosts)
    return {
      states,
      set
    }
  }
}
</script>
<style>
</style>

