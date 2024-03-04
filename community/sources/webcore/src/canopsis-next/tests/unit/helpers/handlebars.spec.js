import { getTemplateVariables } from '@/helpers/handlebars';

describe('getTemplateVariables', () => {
  it('Return empty variables without template', () => {
    expect(getTemplateVariables('')).toEqual([]);
  });

  it('Return variable with simple template', () => {
    expect(getTemplateVariables('{{value}}')).toEqual(['value']);
  });

  it('Return variable with list template', () => {
    expect(
      getTemplateVariables(
        `
          <ul>
            <li><i>{{italic.value}}</i></li>
            <li><b>{{bold.value}}</b></li>
          </ul>
        `,
      ),
    ).toEqual(['italic.value', 'bold.value']);
  });

  it('Return variable with compare helper', () => {
    expect(
      getTemplateVariables(
        `
         <p>{{#compare compare.first '==' compare.second }}</p>
         <p>{{uppercase compare.content}}</p>
         <p>{{/compare}}</p>
        `,
      ).sort(),
    ).toEqual(['compare.first', 'compare.second', 'compare.content'].sort());
  });

  it('Return variable with request helper', () => {
    expect(
      getTemplateVariables(
        `
         <p>{{#request url=(concat "http://example.com/" request.concat.value)}}
          {{uppercase (concat request.uppercase.concat.value1 request.uppercase.concat.value2)}}
          text
         {{/request}}</p>
        `,
      ).sort(),
    ).toEqual(['request.concat.value', 'request.uppercase.concat.value1', 'request.uppercase.concat.value2'].sort());
  });

  it('Return variable with internal-link helper', () => {
    expect(
      getTemplateVariables('{{internal-link href="/admin/maps" text=link.text}}'),
    ).toEqual(['link.text']);
  });

  it('Return variable with copy helper', () => {
    expect(
      getTemplateVariables('{{#copy copy.arg}}{{copy.content}}{{/copy}}').sort(),
    ).toEqual(['copy.arg', 'copy.content'].sort());
  });
});
