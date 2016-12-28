FROM microsoft/dotnet:1.0.0-preview2-sdk

ADD src/Worker /app/src/Worker
WORKDIR /app


RUN dotnet restore -v minimal src/
RUN dotnet publish -c Release -o ./ src/Worker/
RUN rm -rf src/ $HOME/.nuget/

CMD dotnet Worker.dll
