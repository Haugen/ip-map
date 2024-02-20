const formResult = document.querySelector('#form-result');
const output = document.querySelector('#output');
const submitButton = document.querySelector('#submit-button');

submitButton.addEventListener('click', () => {
  formResult.innerHTML = '';
  output.innerHTML = '';
});

const observer = new MutationObserver((mutations) => {
  try {
    if (mutations[0]?.addedNodes.length === 0) {
      return;
    }

    const data = JSON.parse(mutations[0]?.addedNodes[0]?.textContent);
    data.sort((a, b) => a.ResponseTime - b.ResponseTime);

    for (let i = 0; i < data.length; i++) {
      const domNode = document.createElement('div');
      domNode.textContent = `${data[i].Location.Lat}, ${data[i].Location.Lon}, ${data[i].ResponseTime}`;

      if (i === 0) {
        output.prepend(domNode);
      } else {
        // setTimeout not working properly like this! Duh!
        setTimeout(() => {
          output.prepend(domNode);
        }, data[i].ResponseTime - data[i - 1].ResponseTime);
      }
    }
  } catch (error) {
    console.error("You've gone messed up the JSON", error);
  }
});

observer.observe(formResult, {
  childList: true,
});
