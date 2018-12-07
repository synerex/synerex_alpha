// @flow

import { type Icon, Marker as LeafletMarker , DomUtil} from 'leaflet'
import React from 'react'
import { LeafletProvider, withLeaflet, MapLayer } from 'react-leaflet'
import type { LatLng, MapLayerProps } from 'react-leaflet'

type LeafletElement = LeafletMarker
type Props = {
    icon?: Icon,
    draggable?: boolean,
    opacity?: number,
    position: LatLng,
    zIndexOffset?: number,
} & MapLayerProps

(function() {
    // save these original methods before they are overwritten
    var proto_initIcon = LeafletMarker.prototype._initIcon;
    var proto_setPos = LeafletMarker.prototype._setPos;

    var oldIE = (DomUtil.TRANSFORM === 'msTransform');

    LeafletMarker.addInitHook(function () {
        var iconOptions = this.options.icon && this.options.icon.options;
        var iconAnchor = iconOptions && this.options.icon.options.iconAnchor;
        if (iconAnchor) {
            iconAnchor = (iconAnchor[0] + 'px ' + iconAnchor[1] + 'px');
        }
        this.options.rotationOrigin = this.options.rotationOrigin || iconAnchor || 'center bottom' ;
        this.options.rotationAngle = this.options.rotationAngle || 0;

        // Ensure marker keeps rotated during dragging
        this.on('drag', function(e) { e.target._applyRotation(); });
    });

    LeafletMarker.include({
        _initIcon: function() {
            proto_initIcon.call(this);
        },

        _setPos: function (pos) {
            proto_setPos.call(this, pos);
            this._applyRotation();
        },

        _applyRotation: function () {
            if(this.options.rotationAngle) {
                this._icon.style[DomUtil.TRANSFORM+'Origin'] = this.options.rotationOrigin;

                if(oldIE) {
                    // for IE 9, use the 2D rotation
                    this._icon.style[DomUtil.TRANSFORM] = 'rotate(' + this.options.rotationAngle + 'deg)';
                } else {
                    // for modern browsers, prefer the 3D accelerated version
                    this._icon.style[DomUtil.TRANSFORM] += ' rotateZ(' + this.options.rotationAngle + 'deg)';
                }
            }
        },

        setRotationAngle: function(angle) {
            this.options.rotationAngle = angle;
            this.update();
            return this;
        },

        setRotationOrigin: function(origin) {
            this.options.rotationOrigin = origin;
            this.update();
            return this;
        }
    });
})();


class RMarker extends MapLayer<LeafletElement, Props> {


    createLeafletElement(props: Props): LeafletElement {
        const el = new LeafletMarker(props.position, this.getOptions(props))
        this.contextValue = { ...props.leaflet, popupContainer: el }
        return el
    }

    updateLeafletElement(fromProps: Props, toProps: Props) {
        if (toProps.position !== fromProps.position) {
            this.leafletElement.setLatLng(toProps.position)
        }
        if (toProps.icon !== fromProps.icon) {
            this.leafletElement.setIcon(toProps.icon)
        }
        if (toProps.zIndexOffset !== fromProps.zIndexOffset) {
            this.leafletElement.setZIndexOffset(toProps.zIndexOffset)
        }
        if (toProps.opacity !== fromProps.opacity) {
            this.leafletElement.setOpacity(toProps.opacity)
        }
        if (toProps.draggable !== fromProps.draggable) {
            if (toProps.draggable === true) {
                this.leafletElement.dragging.enable()
            } else {
                this.leafletElement.dragging.disable()
            }
        }
        if (toProps.rotationAngle !== fromProps.rotationAngle) {
            this.leafletElement.setRotationAngle(toProps.rotationAngle)
        }
        if (toProps.rotationOrigin !== fromProps.rotationOrigin) {
            this.leafletElement.setRotationOrigin(toProps.rotationOrigin)
        }

    }

    render() {
        const { children } = this.props
        return children == null || this.contextValue == null ? null : (
            <LeafletProvider value={this.contextValue}>{children}</LeafletProvider>
        )
    }
}

export default withLeaflet(RMarker)
