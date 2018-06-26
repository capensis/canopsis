<template>
  <v-container>
    <v-layout>
      <v-flex>
        <v-form>
          <v-text-field v-model="formData.name" label="Name" required></v-text-field>
          <v-text-field v-model="formData.description" label="Description"></v-text-field>
          <v-text-field v-model="formData.value" label="Value"></v-text-field>
          <v-btn @click="submit">
            Submit
          </v-btn>
        </v-form>
      </v-flex>
    </v-layout>
  </v-container>
</template>

<script>
export default {
  props: {
    infoObject: {
      type: Object,
      required: false,
    },
  },
  data() {
    return {
      formData: {
        name: null,
        description: null,
        value: null,
      },
    };
  },
  created() {
    if (this.infoObject) {
      this.formData.name = this.infoObject.name;
      this.formData.description = this.infoObject.description;
      this.formData.value = this.infoObject.value;
    }
  },
  methods: {
    submit() {
      this.$emit('submit', this.formData);
    },
  },
};
</script>
