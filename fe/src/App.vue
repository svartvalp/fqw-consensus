<template>
  <WidgetContainerModal/>
  <nav class="navbar navbar-expand-lg navbar-light bg-light" style="height: 8%;">
    <div class="container-fluid">
      <a class="navbar-brand" href="#">Мониторинг</a>
    </div>
  </nav>
  <div class="container" style="height: 88%; margin-top: 2%">
    <div class="row justify-content-center" style="margin-bottom: 3%;">
      <div class="col-md-2">
        <button type="button" class="btn btn-primary" @click="open()">Отправить состояние</button>
      </div>
    </div>
    <div class="row align-items-center h-75">
      <router-view :states="states">
      </router-view>
    </div>
  </div>
</template>
<script>
import axios from 'axios';
import {reactive, ref} from "vue";
import {container, openModal} from "jenesius-vue-modal"
import ReqModal from "./components/ReqModal";

export default {
  name: 'App',
  components: {
    WidgetContainerModal: container,
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

    const loadData = async (hosts) => {
      setInterval(async function () {
        let newStates = []
        for (let i = 0; i < hosts.length; i++) {
          try {
            let res = await axios.get("http://" + hosts[i] + "/state")
            let state = res.data
            state.host = hosts[i]
            newStates.push(state)
          } catch (e) {
            newStates.push({
              state: 'Dead',
              host: hosts[i]
            })
          }
        }
        set(newStates)
      }, 100)
    }

    loadData(hosts)
    return {
      states,
      set,
      hosts,
    }
  },
  methods: {
    open() {
      openModal(ReqModal, {
        hosts: this.hosts
      })
    }
  }
}
</script>
<style>
.container {
  z-index: -1;
}

.modal-container {
  z-index: 999;
}
</style>

