<template>
  <v-container>
    <v-layout>
      <v-flex>
        <v-form>
          <v-text-field v-model="formData.name" name="name" v-validate="'required|unique-name'" label="Name" required>
          </v-text-field>
          <span class="red--text">{{ errors.first('name') }}</span>
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
  inject: ['$validator'],
  props: {
    infoObject: {
      type: Object,
      required: false,
    },
    forbiddenNames: {
      type: Array,
      default: () => [
        'test',
      ],
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
    this.createUniqueValidationRule();
    if (this.infoObject) {
      this.formData.name = this.infoObject.name;
      this.formData.description = this.infoObject.description;
      this.formData.value = this.infoObject.value;
    }
  },
  methods: {
    submit() {
      this.$validator.validate()
        .then((result) => {
          if (result) {
            this.$emit('submit', this.formData);
          }
        });
    },
    createUniqueValidationRule() {
      this.$validator.extend('unique-name', {
        getMessage: () => this.$t('validator.unique'),
        validate: value => !this.forbiddenNames.includes(value),
      });
    },
  },
};
</script>
