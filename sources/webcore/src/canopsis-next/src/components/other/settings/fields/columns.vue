<template lang="pug">
  v-list-group
    v-list-tile(slot="activator") {{$t('settings.columnNames')}}
    v-container
      v-card
        v-layout.pt-2(justify-space-between)
          v-flex(xs3)
            v-layout.text-xs-center.pl-2(justify-space-between)
              v-flex(xs1)
                v-icon arrow_upward
              v-flex(xs5)
                v-icon arrow_downward
          v-flex.d-flex(xs3)
            div.text-xs-right.pr-2
              v-icon(color="red") close
        v-layout(justify-center wrap)
          v-flex(xs11)
            v-text-field(:placeholder="$t('common.label')")
          v-flex(xs11)
            v-text-field(:placeholder="$t('common.value')")
</template>
