<template lang="pug">
  v-layout(row, align-center)
    c-alarm-tag-chip(
      v-for="tag in inlineTags",
      :key="tag",
      :color="getTagColor(tag)",
      @click="selectTag(tag)"
    ) {{ tag }}
    v-menu(
      v-if="dropDownTags.length",
      bottom,
      left
    )
      template(#activator="{ on }")
        v-btn.ma-1(v-on="on", color="grey", icon, small)
          v-icon(color="white", small) more_horiz
      v-card
        v-card-text
          c-alarm-tag-chip(
            v-for="tag in dropDownTags",
            :key="tag",
            :color="getTagColor(tag)",
            @click="selectTag(tag)"
          ) {{ tag }}
</template>

<script>
import { entitiesAlarmTagMixin } from '@/mixins/entities/alarm-tag';

export default {
  mixins: [entitiesAlarmTagMixin],
  props: {
    alarm: {
      type: Object,
      required: true,
    },
    selectedTag: {
      type: String,
      default: '',
    },
  },
  computed: {
    tags() {
      return [...(this.alarm.tags ?? [])].sort((firstTag, secondTag) => {
        if (firstTag === this.selectedTag) {
          return -1;
        }

        if (secondTag === this.selectedTag) {
          return 0;
        }

        if (firstTag < secondTag) {
          return -1;
        }

        if (firstTag > secondTag) {
          return 1;
        }

        return 0;
      });
    },

    inlineTags() {
      return this.tags.slice(0, 2);
    },

    dropDownTags() {
      return this.tags.slice(2);
    },
  },
  methods: {
    selectTag(tag) {
      this.$emit('select', tag);
    },
  },
};
</script>
