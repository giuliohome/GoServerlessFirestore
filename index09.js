const sendJSON = async () => {
			
			let result = document.querySelector('.result');
			let message = document.querySelector('#name');
			
			// Creating a XHR object
			let xhr = new XMLHttpRequest();
			let url = "/cloudfunction1";
	
			// open a connection
			xhr.open("POST", url, true);
			
			// Set the request header i.e. which type of content you are sending
			xhr.setRequestHeader("Content-Type", "application/json");
			const token = await auth0.getTokenSilently();
			xhr.setRequestHeader("Authorization", `Bearer ${token}`);

			// Create a state change callback
			xhr.onreadystatechange = function () {
				if (xhr.readyState === 4 && xhr.status === 200) {

					// Print received data from server
					result.innerHTML = this.responseText;

				}
			};

			// Converting JSON data to string
			var data = JSON.stringify({ "message": message.value });

			// Sending data with the request
			xhr.withCredentials = true;
			xhr.send(data);
		}
const writtenJSON = async () => {
			
			let result = document.querySelector('.result');
			let message = document.querySelector('#name');
			let written = document.querySelector('#written');
			
			// Creating a XHR object
			let xhr = new XMLHttpRequest();
			let url = "/cloudfunction1";
	
			// open a connection
			xhr.open("POST", url, true);

			// Set the request header i.e. which type of content you are sending
			xhr.setRequestHeader("Content-Type", "application/json");
			const token = await auth0.getTokenSilently();
			xhr.setRequestHeader("Authorization", `Bearer ${token}`);

			// Create a state change callback
			xhr.onreadystatechange = function () {
				if (xhr.readyState === 4 && xhr.status === 200) {

					// Print received data from server
					result.innerHTML = this.responseText;

				}
			};

			// Converting JSON data to string
			var data = JSON.stringify({ "message": message.value, "operation": "write", "written": written.value });

			// Sending data with the request
			xhr.withCredentials = true;
			xhr.send(data);
		}
const queryJSON = async () => {
			
			result = document.getElementById("myTextarea");
			// Creating a XHR object
			let xhr = new XMLHttpRequest();
			let url = "/cloudfunction1";
	
			// open a connection
			xhr.open("POST", url, true);

			// Set the request header i.e. which type of content you are sending
			xhr.setRequestHeader("Content-Type", "application/json");
			const token = await auth0.getTokenSilently();
			xhr.setRequestHeader("Authorization", `Bearer ${token}`);

			// Create a state change callback
			xhr.onreadystatechange = function () {
				if (xhr.readyState === 4 && xhr.status === 200) {

					// Print received data from server
					result.value = this.responseText;

				}
			};

			// Converting JSON data to string
			var data = JSON.stringify({ "operation": "query" });

			// Sending data with the request
			xhr.withCredentials = true;
			xhr.send(data);
}				
