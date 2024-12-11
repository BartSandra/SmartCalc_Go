ifeq ($(shell uname -s), Linux)
OS = linux
else ifeq ($(shell uname -s), Darwin)
OS = darwin
endif

CXXFLAGS = -std=c++11 -fPIC -g
CPPFLAGS = -I/usr/include/c++/13 -I/usr/include/x86_64-linux-gnu/c++/13 -I/usr/include/x86_64-linux-gnu/c++/13/
FYNE = ~/go/bin/fyne
NAME = SmartCalc_v3.0

.PHONY: all install build clean dist uninstall start test style

deps:
	@command -v $(FYNE) >/dev/null 2>&1 || { echo >&2 "Fyne is not installed. Installing..."; go install fyne.io/fyne/v2/cmd/fyne@latest; }

install: uninstall deps build
	@echo "Installing application..."
	# Создаем необходимые каталоги
	@mkdir -p build/Contents/Resources
	# Копируем .env файл, если он существует
	@cp .env build/ || echo ".env file not found, skipping copy"
	# Копируем необходимые файлы из ранее собранной структуры
	@cd cmd && $(FYNE) package -os $(OS) -name SmartCalc_v3.0 -icon ../images/icon.png
	@cd cmd && tar -xf SmartCalc_v3.0.tar.xz
	# Перемещаем исполняемый файл
	@mv cmd/usr/local/bin/cmd build/SmartCalc_v3.0
	# Копируем вспомогательные файлы (иконку и .desktop)
	@cp cmd/usr/local/share/applications/SmartCalc_v3.0.desktop build/Contents/Resources/
	@cp cmd/usr/local/share/pixmaps/SmartCalc_v3.0.png build/Contents/Resources/
	# Копируем дополнительные ресурсы
	@cp internal/model/model/libmodel.so build/Contents/Resources/
	@cp help/help.md build/Contents/Resources/
	@cp history/history.txt build/Contents/Resources/
	@echo "Installation complete. Executable is in build/SmartCalc_v3.0"

dist: deps build
	@echo "Creating distribution package..."
	# Создаем необходимые каталоги
	@mkdir -p build/Contents/Resources
	# Копируем .env файл, если он существует
	@cp .env build/ || echo ".env file not found, skipping copy"
	# Упаковываем приложение с помощью fyne
	@cd cmd && $(FYNE) package -os $(OS) -name SmartCalc_v3.0 -icon ../images/icon.png
	@cd cmd && tar -xf SmartCalc_v3.0.tar.xz
	# Перемещаем исполняемый файл
	@mv cmd/usr/local/bin/cmd build/SmartCalc_v3.0
	# Копируем вспомогательные файлы (иконку и .desktop)
	@cp cmd/usr/local/share/applications/SmartCalc_v3.0.desktop build/Contents/Resources/
	@cp cmd/usr/local/share/pixmaps/SmartCalc_v3.0.png build/Contents/Resources/
	# Копируем дополнительные ресурсы
	@cp internal/model/model/libmodel.so build/Contents/Resources/
	@cp help/help.md build/Contents/Resources/
	@cp history/history.txt build/Contents/Resources/
	# Создаем архив для дистрибутива
	@cd build && tar -czf ~/Desktop/SmartCalc_v3.0.tar.gz SmartCalc_v3.0
	@echo "Distribution package created: ~/Desktop/SmartCalc_v3.0.tar.gz"

build:
	# Сначала компилируем C++ библиотеку
	@g++ $(CXXFLAGS) $(CPPFLAGS) -shared $(PWD)/internal/model/model/model.cc $(PWD)/internal/model/model/credit.cc $(PWD)/internal/model/model/model_wrapper.cc -o $(PWD)/internal/model/model/libmodel.so
	# Затем компилируем Go плагин
	@cd $(PWD)/internal/model/model && go build -buildmode=c-archive -o $(PWD)/internal/model/model/model.a
	@cd $(PWD)/internal/model/model && go build -buildmode=plugin -o $(PWD)/internal/model/model/model.so

start: 
	cd $(PWD)/build && ./SmartCalc_v3.0

test: build
	cd $(PWD)/tests && go test

uninstall: clean
	@echo "Uninstalling..."
	@rm -rf $(PWD)/cmd/usr
	@rm -rf $(PWD)/cmd/Makefile
	@rm -f ~/Desktop/SmartCalc_v3.0.tar.gz
	@rm -f ~/Desktop/SmartCalc_v3.0

clean:
	@rm -f $(PWD)/internal/model/model/libmodel.so
	@rm -f $(PWD)/internal/model/model/model.so
	@rm -f $(PWD)/internal/model/model/model.a
	@rm -rf $(PWD)/build
	@rm -f ~/Desktop/SmartCalc_v3.0.tar.gz
	@rm -rf $(PWD)/cmd/SmartCalc_v3.0.tar.xz

style:
	@echo "Formatting C++ code with clang-format..."
	@find . -name "*.cc" -exec clang-format -i {} \;
	@find . -name "*.h" -exec clang-format -i {} \;
	@echo "C++ code formatting complete."

	@echo "Formatting Go code with gofmt..."
	@gofmt -w .
	@echo "Go code formatting complete."
