<template>
  <v-card>
    <v-layout
      class="pt-2"
      justify-space-between="justify-space-between"
    >
      <v-flex xs5="xs5">
        <span class="pl-2">
          <v-btn
            :disabled="disabledUp"
            icon="icon"
            @click.prevent="$emit('up')"
          >
            <v-icon>arrow_upward</v-icon>
          </v-btn>
          <v-btn
            :disabled="disabledDown"
            icon="icon"
            @click.prevent="$emit('down')"
          >
            <v-icon>arrow_downward</v-icon>
          </v-btn></span>
      </v-flex>
      <v-flex
        class="text-right pr-2"
        xs3="xs3"
      >
        <v-btn
          icon="icon"
          @click.prevent="$emit('remove')"
        >
          <v-icon color="error">
            close
          </v-icon>
        </v-btn>
      </v-flex>
    </v-layout>
    <v-layout
      class="px-3 pb-3"
      justify-center="justify-center"
      column="column"
    >
      <slot />
    </v-layout>
  </v-card>
</template>

<script>
export default {
  props: {
    disabledUp: {
      type: Boolean,
      default: false,
    },
    disabledDown: {
      type: Boolean,
      default: false,
    },
  },
};
</script>
