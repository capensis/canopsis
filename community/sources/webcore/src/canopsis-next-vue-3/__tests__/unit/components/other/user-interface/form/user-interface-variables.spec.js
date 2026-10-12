import { generateShallowRenderer } from '@unit/utils/vue';

import UserInterfaceGeneralForm from '@/components/other/user-interface/form/user-interface-general-form.vue';
import UserInterfaceLoginForm from '@/components/other/user-interface/form/user-interface-login-form.vue';

const { templateVars } = vi.hoisted(() => ({
  templateVars: {
    value: {
      Location: 'string',
      System: {
        ENV_TEST: 'string',
      },
    },
  },
}));

vi.mock('@/hooks/store/modules/template-vars', () => ({
  useTemplateVars: () => ({ templateVars }),
}));

const TextEditorFieldStub = {
  name: 'TextEditorField',
  props: {
    name: {
      type: String,
      default: '',
    },
    variables: {
      type: Array,
      default: () => [],
    },
  },
  template: '<text-editor-field-stub />',
};

const stubs = {
  'c-form-block': true,
  'c-form-block-row': true,
  'c-select': {
    name: 'CSelect',
    props: ['items', 'name'],
    template: '<select class="c-select" />',
  },
  'c-language-field': true,
  'c-theme-field': true,
  'c-timezone-field': true,
  'file-selector': true,
  'text-editor-field': TextEditorFieldStub,
};

const expectedVariables = [
  {
    text: 'env',
    value: '{{ env }}',
    variables: [
      {
        text: 'Location',
        value: '{{ env.Location }}',
      },
      {
        text: 'System',
        value: '{{ env.System }}',
        variables: [
          {
            text: 'ENV_TEST',
            value: '{{ env.System.ENV_TEST }}',
          },
        ],
      },
    ],
  },
];

const selectTextEditorFields = wrapper => wrapper.findAllComponents({ name: 'TextEditorField' });
const selectDateFormatField = wrapper => wrapper.findComponent({ name: 'CSelect' });

describe('user-interface variables', () => {
  test('shows the three localized date format choices', () => {
    const factory = generateShallowRenderer(UserInterfaceGeneralForm, { stubs });
    const wrapper = factory({
      props: {
        form: {},
      },
    });

    expect(selectDateFormatField(wrapper).props('items')).toEqual([
      {
        text: '31/12/2000 00:00:00',
        value: '02/01/2006 15:04:05',
      },
      {
        text: '31/12/2000 (Mon) 00:00:00',
        value: '02/01/2006 (Mon) 15:04:05',
      },
      {
        text: 'Mon 31/12/2000 00:00:00',
        value: 'Mon 02/01/2006 15:04:05',
      },
    ]);
  });

  test('passes environment variables to the app title editor', () => {
    const factory = generateShallowRenderer(UserInterfaceGeneralForm, { stubs });
    const wrapper = factory({
      props: {
        form: {},
      },
    });

    const appTitleField = selectTextEditorFields(wrapper).find(
      field => field.props('name') === 'app_title',
    );

    expect(appTitleField.props('variables')).toEqual(expectedVariables);
  });

  test.each(['footer', 'login_page_description'])(
    'passes environment variables to the %s editor',
    (name) => {
      const factory = generateShallowRenderer(UserInterfaceLoginForm, { stubs });
      const wrapper = factory({
        props: {
          form: {},
        },
      });

      const field = selectTextEditorFields(wrapper).find(editor => editor.props('name') === name);

      expect(field.props('variables')).toEqual(expectedVariables);
    },
  );
});
