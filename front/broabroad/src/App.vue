<template>
  <div id="app" class="container">
    <div class="row">
      <div class="col-md-6 offset-md-3 py-5">
        <h1>Broabroad</h1>

        <form v-on:submit.prevent="makeSeekRequest">
          <div class="form-group">
            <input v-model="websiteUrl" type="text" id="website-input" placeholder="Enter a website" class="form-control">
          </div>
          <div class="form-group">
            <button class="btn btn-primary">Create a seek request</button>
          </div>
        </form>
      </div>
    </div>
  </div>
</template>

<script>
import axios from 'axios';

export default {
  name: 'App',

  data() { return {
    websiteUrl: '',
  } },

  methods: {
    makeSeekRequest() {
      var username = 'test';
      var password = 'qwerty1';
      var basicAuth = 'Basic ' + btoa(username + ':' + password);
      console.log(basicAuth);
      axios.post("http://localhost:9090/seek_requests", {}, {
        headers: { 'Authorization': basicAuth }
      }).then((response) => {
        console.log(response);
      }).catch((error) => {
        window.alert(`The API returned an error: ${error}`);
      })
    }
  }
}
</script>