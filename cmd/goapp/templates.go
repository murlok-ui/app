// Code generated by go generate; DO NOT EDIT.
package main

const appxManifestTmpl = `
<?xml version="1.0" encoding="utf-8"?>
<Package xmlns="http://schemas.microsoft.com/appx/manifest/foundation/windows10" 
    xmlns:mp="http://schemas.microsoft.com/appx/2014/phone/manifest" 
    xmlns:uap="http://schemas.microsoft.com/appx/manifest/uap/windows10" 
    xmlns:rescap="http://schemas.microsoft.com/appx/manifest/foundation/windows10/restrictedcapabilities" 
    xmlns:desktop="http://schemas.microsoft.com/appx/manifest/desktop/windows10" IgnorableNamespaces="uap mp rescap desktop">
    <Identity Name="4a99c091-2185-434d-8088-4b6c643f38fc" Publisher="CN=Maxence" Version="1.0.0.0" />
    <mp:PhoneIdentity PhoneProductId="4a99c091-2185-434d-8088-4b6c643f38fc" PhonePublisherId="00000000-0000-0000-0000-000000000000" />
    <Properties>
        <DisplayName>{{.Name}}</DisplayName>
        <PublisherDisplayName>{{.Publisher}}</PublisherDisplayName>
        <Logo>Assets\StoreLogo.png</Logo>
    </Properties>
    <Dependencies>
        <TargetDeviceFamily Name="Windows.Universal" MinVersion="10.0.14393.0" MaxVersionTested="10.0.16299.15" />
    </Dependencies>
    <Resources>
        <Resource Language="en-us" />
        <Resource uap:Scale="100" />
        <Resource uap:Scale="125" />
        <Resource uap:Scale="150" />
        <Resource uap:Scale="200" />
        <Resource uap:Scale="400" />
    </Resources>
    <Applications>
        <Application Id="{{.ID}}" Executable="uwp.exe" EntryPoint="uwp.App">
            <uap:VisualElements DisplayName="{{.Name}}" Description="{{.Description}}" BackgroundColor="transparent" Square150x150Logo="Assets\Square150x150Logo.png" Square44x44Logo="Assets\Square44x44Logo.png">
                <uap:DefaultTile Wide310x150Logo="Assets\Wide310x150Logo.png" Square310x310Logo="Assets\Square310x310Logo.png" Square71x71Logo="Assets\Square71x71Logo.png">
                    <uap:ShowNameOnTiles>
                        <uap:ShowOn Tile="square150x150Logo" />
                        <uap:ShowOn Tile="wide310x150Logo" />
                        <uap:ShowOn Tile="square310x310Logo" />
                    </uap:ShowNameOnTiles>
                </uap:DefaultTile>
            </uap:VisualElements>
            <Extensions>
                <desktop:Extension Category="windows.fullTrustProcess" Executable="{{.Executable}}" />
                <uap:Extension Category="windows.appService">
                    <uap:AppService Name="goapp" />
                </uap:Extension>
                <uap5:Extension Category="windows.appExecutionAlias" Executable="uwp.exe" EntryPoint="uwp.App}">
                    <uap5:AppExecutionAlias>
                        <uap5:ExecutionAlias Alias="{{.Executable}}" />
                    </uap5:AppExecutionAlias>
                </uap5:Extension>
            </Extensions>
        </Application>
    </Applications>
    <Capabilities>
        <Capability Name="internetClient" />
        <rescap:Capability Name="runFullTrust"/>
    </Capabilities>
</Package>`

