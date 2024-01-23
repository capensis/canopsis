<template>
  <widget-settings-flat-item :title="title">
    <template #actions>
      <v-layout justify-end>
        <v-btn
          v-if="isEmpty && addable"
          class="primary"
          small
          @click="$emit('create', $event)"
        >
          <span>{{ $t('common.create') }}</span>
        </v-btn>
        <template v-else>
          <v-btn
            class="primary"
            small
            @click="$emit('edit', $event)"
          >
            <span>{{ $t('common.edit') }}</span>
          </v-btn>
          <v-btn
            v-if="removable"
            class="error"
            small
            @click="$emit('delete', $event)"
          >
            <v-icon>delete</v-icon>
          </v-btn>
        </template>
      </v-layout>
    </template>
  </widget-settings-flat-item>
</template>

<script>
import WidgetSettingsFlatItem from '@/components/sidebars/partials/widget-settings-flat-item.vue';

export default {
  components: { WidgetSettingsFlatItem },
  props: {
    title: {
      type: String,
      required: true,
    },
    isEmpty: {
      type: Boolean,
      default: true,
    },
    addable: {
      type: Boolean,
      default: false,
    },
    removable: {
      type: Boolean,
      default: false,
    },
  },
};
</script>
