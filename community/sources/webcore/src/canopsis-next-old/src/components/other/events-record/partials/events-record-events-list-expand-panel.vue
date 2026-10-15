<template>
  <div class="secondary pa-3">
    <v-card>
      <v-card-text>
        <c-json-treeview :json="preparedEvent | json" />
      </v-card-text>
    </v-card>
  </div>
</template>
<script>
import { computed } from 'vue';

export default {
  props: {
    event: {
      type: Object,
      required: true,
    },
  },
  setup(props) {
    const preparedEvent = computed(() => {
      let decodedEvent = props.event.event;

      try {
        decodedEvent = atob(decodedEvent);
        decodedEvent = JSON.parse(decodedEvent);
      } catch (err) {
        console.warn(err);
      }

      return {
        ...props.event,

        event: decodedEvent,
      };
    });

    return {
      preparedEvent,
    };
  },
};
</script>
