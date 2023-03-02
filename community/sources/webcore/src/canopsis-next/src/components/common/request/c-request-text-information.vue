<template lang="pug">
  v-layout(column)
    request-information-content(:text="data.request")
    request-information-content.mt-4(v-if="data.response", :text="data.response")
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
