class LoginPage {
  constructor(application) {
    this.application = application;
  }

  async login() {
    this.application.navigate('/login');

    await this.application.waitElement('form');

    await this.application.page.type('input[name="username"]', 'root');
    await this.application.page.type('input[name="password"]', 'root');
    await this.application.page.click('button[type="submit"]');

    await this.application.waitElement('nav');
  }
}

module.exports = {
  LoginPage,
};
