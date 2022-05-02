<template>
  <div class="card modl">
    <div class="card-body">
      <select class="form-select" required v-model="host">
        <option v-for="host in this.hosts" :key="host">{{host}}</option>
      </select>
      <label for="keyIn" class="form-label" >Ключ</label>
      <input class="form-control" id="keyIn" v-model="key">
      <label for="valIn" class="form-label">Значение</label>
      <input class="form-control" id="valIn" v-model="val">
      <button type="button" class="btn btn-primary" @click="send()" >Отправить</button>
    </div>
  </div>
</template>
<script>
import {ref} from "vue";
import axios from "axios";

export default {
  name: 'ReqModal',
  props: {
    hosts: Array,
  },
  setup(){
    const key = ref("")
    const val = ref("")
    const host = ref("")
    return{
      key,
      val,
      host
    }
  },
  methods: {
    async send() {
      try {
        let res = await axios.post("http://" + this.host + "/store", {
          key: this.key,
          value: this.val
        })
        console.log(res)
      }
      catch (e) {
        console.log(e)
      }
    }
  }
}
</script>