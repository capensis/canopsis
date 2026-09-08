<template>
  <div class="bg-secondary pa-3">
    <v-card>
      <v-card-text>
        <c-object-treeview
          v-model:opened="openedItems"
          :root-name="rootName"
          :value="preparedEvent"
          :value-formatter="formatValue"
          value-max-width="720px"
          strong-title
        />
      </v-card-text>
    </v-card>
  </div>
</template>

<script>
import { computed, ref } from 'vue';

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

    const rootName = 'event';
    const openedItems = ref([rootName]);
    const formatValue = value => String(value);

    return {
      rootName,
      preparedEvent,
      openedItems,
      formatValue,
    };
  },
};
</script>
