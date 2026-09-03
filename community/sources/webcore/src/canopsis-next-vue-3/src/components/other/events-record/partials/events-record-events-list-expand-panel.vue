<template>
  <div class="bg-secondary pa-3">
    <v-card>
      <v-card-text>
        <c-treeview
          v-model:opened="openedItems"
          :items="eventTree"
          item-key="path"
          item-title="name"
        >
          <template #prepend="{ item }">
            <div v-if="item.isArray" class="text-caption font-italic text-medium-emphasis">
              {{ `(${$t('common.variableTypes.array')})` }}
            </div>
          </template>

          <template #label="{ item, leaf }">
            <div class="events-record-event-field">
              <strong>{{ item.name }}</strong>
              <span v-if="leaf" class="pl-1">:
                <c-ellipsis
                  :text="formatValue(item.value)"
                  class="
                    pl-1
                    d-inline-block
                    text-medium-emphasis
                    text-body-1
                    pre-wrap
                    events-record-event-field__value
                  "
                />
              </span>
              <span v-else-if="!item.children?.length" class="pl-1">:
                <div class="pl-1 d-inline-block text-medium-emphasis text-body-1 font-italic">
                  {{ $t('common.emptyObject') }}
                </div>
              </span>
            </div>
          </template>
        </c-treeview>
      </v-card-text>
    </v-card>
  </div>
</template>

<script>
import { computed, ref, watch } from 'vue';

import { convertObjectToTreeview } from '@/helpers/treeview';

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

    const eventTree = computed(() => [convertObjectToTreeview(preparedEvent.value, 'event')]);
    const openedItems = ref([]);
    const formatValue = value => String(value);

    watch(
      eventTree,
      (tree) => {
        openedItems.value = tree.map(item => item.path);
      },
      { immediate: true },
    );

    return {
      eventTree,
      openedItems,
      formatValue,
    };
  },
};
</script>

<style lang="scss" scoped>
.events-record-event-field {
  min-width: 0;

  &__value {
    max-width: 720px;
    vertical-align: bottom;
  }
}
</style>
