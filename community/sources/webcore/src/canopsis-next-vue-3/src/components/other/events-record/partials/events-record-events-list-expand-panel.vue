<template>
  <div class="bg-secondary pa-3">
    <v-card>
      <v-card-text>
        <c-json-treeview :json="preparedEventString" />
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

    const preparedEventString = computed(() => JSON.stringify(preparedEvent.value));

    return {
      preparedEventString,
    };
  },
};
</script>
