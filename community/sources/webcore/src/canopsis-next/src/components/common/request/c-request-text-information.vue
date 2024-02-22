<template>
  <v-layout column>
    <request-information-content :text="data.request" />
    <request-information-content
      v-if="data.response"
      :text="data.response"
      class="mt-4"
    />
  </v-layout>
</template>

<script>
import RequestInformationContent from '@/components/common/request/partials/request-information-content.vue';

export default {
  components: { RequestInformationContent },
  props: {
    value: {
      type: String,
      required: true,
    },
  },
  computed: {
    data() {
      const [request, response] = this.value.split('\n\n');

      return {
        request,
        response,
      };
    },
  },
};
</script>
