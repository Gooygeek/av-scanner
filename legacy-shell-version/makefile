
.PHONY: install_service
install_service:
	sudo cp ./systemd-av-watcher.service /etc/systemd/system/av-watcher.service

.PHONY: start_service
start_service:
	sudo systemctl restart av-watcher.service
