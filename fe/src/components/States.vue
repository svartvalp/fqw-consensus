<template>
  <div>
    <div class="row align-items-center h-100">
      <State v-for="state in states" :state="state" :key="state"></State>
    </div>
  </div>
</template>
<script>
import State from "./State";
export default {
  name:'States',
  props: {
    states: Array,
  },
  components: {
    State
  },
  setup(){

  }
}
</script>